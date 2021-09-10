package v1alpha1

type EventAttribute string

const (
	Action        EventAttribute = "action"
	Reference     EventAttribute = "reference"
	Repo          EventAttribute = "repository"
	Sender        EventAttribute = "sender"
	Number        EventAttribute = "number"
	Branch        EventAttribute = "branch"
	SourceBranch  EventAttribute = "sourcebranch"
	TargetBranch  EventAttribute = "targetbranch"
	Commit        EventAttribute = "commit"
	CommitMessage EventAttribute = "commitmessage"
	Tag           EventAttribute = "tag"
)
