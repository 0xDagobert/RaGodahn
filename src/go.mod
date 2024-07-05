module RaGodahn

go 1.22.5

require (
	RaGodahn/pslistwin v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.20.0
)

replace RaGodahn/pslistwin => ./PsListWin

replace RaGodahn/injector => ./injector

replace RaGodahn/packer => ./packer

replace RaGodahn/payload => ./payload
