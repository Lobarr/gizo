package job

//EnvironmentVariable stores key and value of env variables
type EnvironmentVariable struct {
	Key   string
	Value string
}

//NewEnv initializes an environment variable
func NewEnv(key, value string) *EnvironmentVariable {
	return &EnvironmentVariable{
		Key:   key,
		Value: value,
	}
}

//GetKey returns key
func (env EnvironmentVariable) GetKey() string {
	return env.Key
}

//GetValue returns value
func (env EnvironmentVariable) GetValue() string {
	return env.Value
}

//EnvironmentVariables list of environment variables
type EnvironmentVariables []EnvironmentVariable

//NewEnvVariables initializes environment variables
func NewEnvVariables(variables ...EnvironmentVariable) EnvironmentVariables {
	return variables
}
