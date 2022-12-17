package http

type Error interface {
	HttpCode() int
}
