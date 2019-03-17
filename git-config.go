package gitconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/jiangxin/gitconfig/goconfig"
)

const maxIncludeDepth = 10

// Scope is used to mark where config variable comes from
type Scope uint16

// Define scopes for config variables
const (
	ScopeInclude Scope = 1 << iota
	ScopeSystem
	ScopeGlobal
	ScopeSelf

	ScopeAll  Scope = 0xFFFF
	ScopeMask Scope = ^ScopeInclude
)

// String show user friendly display of scope
func (v *Scope) String() string {
	inc := ""
	if (*v & ScopeInclude) == ScopeInclude {
		inc = "-inc"
	}

	if (*v & ScopeSystem) == ScopeSystem {
		return "system" + inc
	} else if (*v & ScopeGlobal) == ScopeGlobal {
		return "global" + inc
	} else if (*v & ScopeSelf) == ScopeSelf {
		return "self" + inc
	}
	return "unknown" + inc
}

// GitConfig maps section to key-value pairs
type GitConfig map[string]GitConfigKeyValues

// GitConfigKeyValues maps key to values
type GitConfigKeyValues map[string][]GitConfigValue

// GitConfigValue holds value and its scope
type GitConfigValue struct {
	scope Scope
	value string
}

// Keys returns sorted kesy in one section
func (v GitConfigKeyValues) Keys() []string {
	keys := reflect.ValueOf(v).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	sort.Strings(strkeys)
	return strkeys
}

// Set is used to set value
func (v *GitConfigValue) Set(value string) {
	v.value = value
}

// Value is used to show value
func (v GitConfigValue) Value() string {
	return v.value
}

// Scope is used to show user friendly scope
func (v GitConfigValue) Scope() string {
	return v.scope.String()
}

// NewGitConfig returns GitConfig with initialized maps
func NewGitConfig() GitConfig {
	c := make(GitConfig)
	return c
}

// Sections returns sorted sections
func (v GitConfig) Sections() []string {
	keys := reflect.ValueOf(v).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	sort.Strings(strkeys)
	return strkeys
}

// Keys returns all config variable keys (in lower case)
func (v GitConfig) Keys() []string {
	allKeys := []string{}
	for s, keys := range v {
		for key := range keys {
			allKeys = append(allKeys, s+"."+key)
		}
	}
	sort.Strings(allKeys)
	return allKeys
}

// Set will replace old config variable
func (v GitConfig) Set(key, value string) {
	s, k := toSectionKey(key)
	keys := v[s]
	if keys == nil {
		v._add(s, k, value)
		return
	}

	if keys[k] == nil || len(keys[k]) == 0 {
		v._add(s, k, value)
		return
	}

	found := false
	for i := len(keys[k]) - 1; i >= 0; i-- {
		if keys[k][i].scope == ScopeSelf {
			found = true
			keys[k][i].value = value
			break
		}
	}

	if !found {
		keys[k] = append(keys[k],
			GitConfigValue{
				scope: ScopeSelf,
				value: value,
			})
	}
}

// Unset will remove latest setting of a config variable
func (v GitConfig) Unset(key string) {
	v.unset(key, false)
}

// UnsetAll will remove all settings of a config variable
func (v GitConfig) UnsetAll(key string) {
	v.unset(key, true)
}

func (v GitConfig) unset(key string, all bool) {
	s, k := toSectionKey(key)
	keys := v[s]
	if keys == nil {
		return
	}

	if keys[k] == nil || len(keys[k]) == 0 {
		return
	}

	for i := len(keys[k]) - 1; i >= 0; i-- {
		if keys[k][i].scope == ScopeSelf {
			keys[k] = append(keys[k][:i], keys[k][i+1:]...)
			if !all {
				break
			}
		}
	}
}

// Add will add user input key-value pair
func (v GitConfig) Add(key string, value ...string) {
	s, k := toSectionKey(key)
	v._add(s, k, value...)
}

// _add key/value to config variables
func (v GitConfig) _add(section, key string, value ...string) {
	// section, and key are always in lower case
	if _, ok := v[section]; !ok {
		v[section] = make(GitConfigKeyValues)
	}

	if _, ok := v[section][key]; !ok {
		v[section][key] = []GitConfigValue{}
	}
	for _, val := range value {
		v[section][key] = append(v[section][key],
			GitConfigValue{
				scope: ScopeSelf,
				value: val,
			})
	}
}

// Get value from key
func (v GitConfig) Get(key string) string {
	values := v.GetAll(key)
	if values == nil || len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

// GetBool gets boolean from key with default value
func (v GitConfig) GetBool(key string, defaultValue bool) (bool, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	switch strings.ToLower(value) {
	case "yes", "true", "on":
		return true, nil
	case "no", "false", "off":
		return false, nil
	}
	return false, ErrNotBoolValue
}

// GetInt return integer value of key with default
func (v GitConfig) GetInt(key string, defaultValue int) (int, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

// GetInt64 return int64 value of key with default
func (v GitConfig) GetInt64(key string, defaultValue int64) (int64, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseInt(value, 10, 64)
}

// GetUint64 return uint64 value of key with default
func (v GitConfig) GetUint64(key string, defaultValue uint64) (uint64, error) {
	value := v.Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseUint(value, 10, 64)
}

// GetAll gets all values of a key
func (v GitConfig) GetAll(key string) []string {
	section, key := toSectionKey(key)

	values := []string{}

	if v[section] != nil && v[section][key] != nil {
		for _, value := range v[section][key] {
			values = append(values, value.value)
		}
		return values
	}
	return nil
}

// GetRaw gets all values of a key
func (v GitConfig) GetRaw(key string) []GitConfigValue {
	section, key := toSectionKey(key)

	if v[section] != nil && v[section][key] != nil {
		return v[section][key]
	}
	return nil
}

func dequoteKey(key string) string {
	if !strings.ContainsAny(key, "\"'") {
		return key
	}

	keys := []string{}
	for _, k := range strings.Split(key, ".") {
		keys = append(keys, strings.Trim(k, "\"'"))

	}
	return strings.Join(keys, ".")
}

// splitKey will split git config variable to section name and key
func toSectionKey(name string) (string, string) {
	name = strings.ToLower(dequoteKey(name))
	items := strings.Split(name, ".")

	if len(items) < 2 {
		return "", ""
	}
	key := items[len(items)-1]
	section := strings.Join(items[0:len(items)-1], ".")
	return section, key
}

// Parse takes given bytes as configuration file (according to gitconfig syntax)
func Parse(bytes []byte, filename string) (GitConfig, uint, error) {
	var (
		gitCfg = NewGitConfig()
		line   uint
		err    error
		depth  uint
	)

	for {
		var (
			cfg         = NewGitConfig()
			file        string
			gocfg       map[string][]string
			includePath string
		)

		gocfg, line, err = goconfig.Parse(bytes)
		for key, val := range gocfg {
			cfg.Add(key, val...)
		}
		if depth == 0 {
			gitCfg = cfg
		} else {
			gitCfg = gitCfg.Merge(cfg, ScopeInclude)
		}
		includePath = cfg.Get("include.path")
		if includePath == "" {
			break
		}
		file, err = AbsJoin(path.Dir(filename), includePath)
		if err != nil {
			break
		}
		depth++
		// Check circular includes
		if depth >= maxIncludeDepth {
			err = fmt.Errorf("exceeded maximum include depth (%d) while including\n"+
				"\t%s\n"+
				"from"+
				"\t%s\n"+
				"This might be due to circular includes\n",
				maxIncludeDepth,
				filename,
				file)
			break
		}
		filename = file
		bytes, err = ioutil.ReadFile(file)
		if err != nil {
			break
		}
	}
	return gitCfg, line, err
}

// Merge will merge another GitConfig, and new value(s) of the same key will
// append to the end of value list, and new value has higher priority.
func (v GitConfig) Merge(c GitConfig, scope Scope) GitConfig {
	for sec, keys := range c {
		if _, ok := v[sec]; !ok {
			v[sec] = make(GitConfigKeyValues)
		}
		for key, values := range keys {
			if v[sec][key] == nil {
				v[sec][key] = []GitConfigValue{}
			}
			for _, value := range values {
				v[sec][key] = append(v[sec][key],
					GitConfigValue{
						scope: (value.scope & ^ScopeMask) | scope,
						value: value.Value(),
					})

			}
		}
	}
	return v
}

// String returns content of GitConfig ready to save config file
func (v GitConfig) String() string {
	return v.StringOfScope(ScopeSelf)
}

// StringOfScope returns contents with matching scope ready to save config file
func (v GitConfig) StringOfScope(scope Scope) string {
	lines := []string{}
	showInc := false
	if (scope & ScopeInclude) != 0 {
		showInc = true
		scope &= (^ScopeInclude)
	}

	for _, s := range v.Sections() {
		secs := strings.SplitN(s, ".", 2)
		sec := s
		if len(secs) == 2 {
			sec = fmt.Sprintf("%s \"%s\"", secs[0], secs[1])
		}
		once := true
		for _, k := range v[s].Keys() {
			for _, value := range v[s][k] {
				if !showInc && ((value.scope & ScopeInclude) != 0) {
					continue
				}
				if (value.scope & (^ScopeInclude) & scope) == 0 {
					continue
				}

				if once {
					once = false
					lines = append(lines, "["+sec+"]")
				}
				line := "\t" + k + " = "
				quote := false
				if isspace(value.value[0]) || isspace(value.value[len(value.value)-1]) {
					quote = true
				}
				if quote {
					line += "\""
				}
				for _, c := range value.value {
					switch c {
					case '\n':
						line += "\\n"
						continue
					case '\t':
						line += "\\t"
						continue
					case '\b':
						line += "\\b"
						continue
					case '\\':
						line += "\\"
					case '"':
						line += "\\"
					}
					line += string(c)
				}
				if quote {
					line += "\""
				}

				lines = append(lines, line)
			}
		}

	}
	return strings.Join(lines, "\n") + "\n"
}

func isspace(c byte) bool {
	return c == '\t' || c == ' ' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// Save will save git config to file
func (v GitConfig) Save(file string) error {
	if file == "" {
		return fmt.Errorf("cannot save config, unknown filename")
	}

	lockFile := file + ".lock"

	err := ioutil.WriteFile(lockFile, []byte(v.String()), 0644)
	defer os.Remove(lockFile)

	if err != nil {
		return err
	}

	_, err = LoadFile(lockFile, false)
	if err != nil {
		return fmt.Errorf("fail to save '%s': %s", file, err)
	}

	return os.Rename(lockFile, file)
}
