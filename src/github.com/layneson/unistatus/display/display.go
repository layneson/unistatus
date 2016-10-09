package display

//A Status represents an animation or still image which is meant to convey information to the user.
type Status interface {
	//Init is called before the Status is rendered.
	//All initialization logic for the Status should occur within this method.
	Init() error

	//Display is called when the Status must be displayed.
	//Animation must occur within this method, with time.Sleep() called between each frame.
	//The method should run for the given amount of seconds.
	//Returns an error if something doesn't go as intended.
	//Should return nil when the Status is finished being displayed.
	Display(seconds int) error
}
