module RaGodahn

go 1.22.5

require (
	RaGodahn/injector v0.0.0-00010101000000-000000000000
	RaGodahn/pslistwin v0.0.0-00010101000000-000000000000
)

replace RaGodahn/pslistwin => ./PsListWin

replace RaGodahn/injector => ./injector

replace RaGodahn/packer => ./packer

replace RaGodahn/payload => ./payload
