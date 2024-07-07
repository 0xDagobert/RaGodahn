module RaGodahn

go 1.22.5

require (
	RaGodahn/injector v0.0.0-00010101000000-000000000000
	RaGodahn/process_a_injector v0.0.0-00010101000000-000000000000
	RaGodahn/pslistwin v0.0.0-00010101000000-000000000000
)

require (
	RaGodahn/pe v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.22.0
)

replace RaGodahn/pslistwin => ./psListWin

replace RaGodahn/injector => ./injector

replace RaGodahn/process_a_injector => ./processA

replace RaGodahn/pe => ./peInjector
