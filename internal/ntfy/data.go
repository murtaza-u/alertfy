package ntfy

// Data contains all the details of the notification.
type Data struct {
	URL         string
	Title       string
	Description string
	Priority    string
	Tags        string
}

// defaultPriority is the priority level used when none is specified.
const defaultPriority = "default"
