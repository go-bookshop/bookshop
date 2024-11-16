package mock

const (
	UnexpectedUserId = -321321
	CorruptedUserId  = 17364
)

const (
	ValidEmail                     = "valid@mail.com"
	ConflictEmail                  = "conflict@mail.com"
	UnexpectedEmail                = "unexpected@mail.com"
	CorruptedEmail                 = "corrupted@mail.com"
	DuplicateEmail                 = "duplicate@mail.com"
	PanicEmail                     = "panic@mail.com"
	FailedToSendEmail              = "failsend@mail.com"
	NotFoundEmail                  = "notfound@mail.com"
	ActivatedEmail                 = "activated@mail.com"
	SimulateFailTokenCreationEmail = "tokenfail@mail.com"
)

const (
	ValidToken                           = "35FM5TZVYCDCPQBCOWHVOVGE6Q"
	SimulateConflictWriteToken           = "99FM5TZVYCDCPQBCDWHVOVGE6Q"
	SimulateUnexpectedWriteToken         = "99FM5TZVYCDCPQBCDWHVOVGE6B"
	ExpiredToken                         = "IS5BCEAMZJ7X5YD54Y6TW4VA2U"
	SimulateFailToDeleteAllByUserIdToken = "AS5BCEAMZJ7X5YD54Y6TW4VA2U"
)
