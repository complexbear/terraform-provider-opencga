package opencga

// This module contains structs to represent the data returned from OpenCGA API calls

/*
File represents a catalog entry with meta data on a file in a mounted filesystem
*/
type File struct {
	Id        int    `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Type      string `mapstructure:"type"`
	Format    string `mapstructure:"format"`
	Bioformat string `mapstructure:"bioformat"`
	Uri       string `mapstructure:"uri"`
	Path      string `mapstructure:"path"`
}

/*
Project represents the essential data we need to know about an opencga project
This data may be required when creating studies, eg id

This is not an exhaustive set of attributes that make up the opencga project
data type. Add more as and when they are needed by the provider.
*/
type Organism struct {
	ScientificName string `mapstructure:"scientificName"`
	TaxonomyCode   int    `mapstructure:"taxonomyCode"`
	Assembly       string `mapstructure:"assembly"`
}
type Project struct {
	Id          int      `mapstructure:"id"`
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Alias       string   `mapstructure:"alias"`
	Organism    Organism `mapstructure:"organism"`
}

/*
Study represents the essential data we need to know about an opencga study.

This is not an exhaustive set of attributes that make up the opencga study
data type. Add more as and when they are needed by the provider.
*/
type Study struct {
	Id          int    `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Alias       string `mapstructure:"alias"`
	Description string `mapstructure:"description"`
}

/*
To manage user and group permissions for a study we use the ACL
This allows updates to the ACL that is applied to a study
*/
type StudyACL struct {
	Member      string   `mapstructure:"member"`
	Permissions []string `mapstructure:"permissions"`
}

/*
User groups configured for a study, eg @members or @admins
To be used for creating AD groups
*/
type StudyGroup struct {
	Id   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

/*
VariableSet definition, but does not include definition of variables.
Indend to support these via raw json structures in terraform projects.
*/
type VariableSet struct {
	Id          int           `mapstructure:"id"`
	Name        string        `mapstructure:"name"`
	Description string        `mapstructure:"description"`
	Unique      bool          `mapstructure:"unique"`
	Variables   []interface{} `mapstructure:"variables"`
}

/*
Login represents the data returned from a user login request
*/
type Login struct {
	Token string `mapstructure:"token"`
}

/*
Common data struct for all responses from OpenCGA.

Use the ResultType value to determine which struct type to map the
contents of the Result value to.

For example a value of ResultType may be: "org.opencb.opencga.core.models.Project"
The value in Result should be mapped to the Project struct.
*/
type Response struct {
	Id              string        `mapstructure:"id"`
	DbTime          int           `mapstructure:"dbTime"`
	NumResults      int           `mapstructure:"numResults"`
	NumTotalResults int           `mapstructure:"numTotalResults"`
	WarningMsg      string        `mapstructure:"warningMsg"`
	ErrorMsg        string        `mapstructure:"errorMsg"`
	ResultType      string        `mapstructure:"resultType"`
	Results         []interface{} `mapstructure:"result"`
}

/*
Top level struct to represent responses from OpenCGA.

Users should take note of the Error entry which can be used to
determine whether the API call was successful or not.
*/
type ApiResponse struct {
	Error     string     `mapstructure:"error"`
	Responses []Response `mapstructure:"response"`
}
