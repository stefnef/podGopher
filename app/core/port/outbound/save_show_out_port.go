package outbound

type SaveShowPort interface {
	SaveShow(title string) (err error)
	ExistsByTitle(title string) bool
}
