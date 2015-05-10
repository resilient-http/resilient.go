package middlewares

type Wrapper interface {
	SetServers([]string)
	Options()
}

func Wrap(interface{}) {

}
