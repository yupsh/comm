package command

type SuppressColumn1Flag bool

const (
	SuppressColumn1   SuppressColumn1Flag = true
	NoSuppressColumn1 SuppressColumn1Flag = false
)

type SuppressColumn2Flag bool

const (
	SuppressColumn2   SuppressColumn2Flag = true
	NoSuppressColumn2 SuppressColumn2Flag = false
)

type SuppressColumn3Flag bool

const (
	SuppressColumn3   SuppressColumn3Flag = true
	NoSuppressColumn3 SuppressColumn3Flag = false
)

type CheckOrderFlag bool

const (
	CheckOrder   CheckOrderFlag = true
	NoCheckOrder CheckOrderFlag = false
)

type TotalFlag bool

const (
	Total   TotalFlag = true
	NoTotal TotalFlag = false
)

type flags struct {
	SuppressColumn1 SuppressColumn1Flag
	SuppressColumn2 SuppressColumn2Flag
	SuppressColumn3 SuppressColumn3Flag
	CheckOrder      CheckOrderFlag
	Total           TotalFlag
}

func (s SuppressColumn1Flag) Configure(flags *flags) { flags.SuppressColumn1 = s }
func (s SuppressColumn2Flag) Configure(flags *flags) { flags.SuppressColumn2 = s }
func (s SuppressColumn3Flag) Configure(flags *flags) { flags.SuppressColumn3 = s }
func (c CheckOrderFlag) Configure(flags *flags)      { flags.CheckOrder = c }
func (t TotalFlag) Configure(flags *flags)           { flags.Total = t }
