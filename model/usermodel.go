package model

//UserDataStore defines the User type data operations.
type UserDataStore interface {
	GetAll() ([]string, error)
	Create(User) bool
	Edit(User) bool
	Delete(int) bool
}

//User is an instance of an employee in a company.
type User struct {
	FirstName    string
	LastName     string
	Email        string
	Organization string
}

//FileUserModel is an implementation of UserDataStore using the filesystem as a pseudo-database.
type FileUserModel struct {
	Filepath *string
}

//GetAll retrieves all saved users.
func (m FileUserModel) GetAll() ([]string, error) {
	//content, err := ioutil.ReadFile(m.Filepath)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	return []string{}, nil
}

//Create creates a new user and saves it to the "database" file.
func (m FileUserModel) Create(u User) bool {
	return false
}

//Edit modifies the properties of the given user based on UI input.
func (m FileUserModel) Edit(u User) bool {
	return false
}

//Delete finds the specified user by ID and deletes them.
func (m FileUserModel) Delete(id int) bool {
	return false
}
