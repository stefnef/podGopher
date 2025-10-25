package outbound

type SaveShowPort interface {
	SaveShow(title string) (id string, err error)
	ExistsByTitle(title string) bool
}
