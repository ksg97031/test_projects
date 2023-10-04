package runner

/**
  @author: yhy
  @since: 2022/12/13
  @desc: //TODO
**/

var ConfigFileName = "config.yaml"

// Default configuration file,  todo Note: Codeql Do not support the specified folder to run the rules
var defaultYamlByte = []byte(`
python_ql:
  - python/ql/src/Security/ksg97031/Format.ql
`)
