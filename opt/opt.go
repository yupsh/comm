package opt

// Boolean flag types with constants
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

// Flags represents the configuration options for the comm command
type Flags struct {
	SuppressColumn1 SuppressColumn1Flag // Suppress column 1 (lines unique to FILE1) (-1)
	SuppressColumn2 SuppressColumn2Flag // Suppress column 2 (lines unique to FILE2) (-2)
	SuppressColumn3 SuppressColumn3Flag // Suppress column 3 (lines common to both) (-3)
	CheckOrder      CheckOrderFlag      // Check that inputs are sorted (--check-order)
	Total           TotalFlag           // Output summary counts (--total)
}

// Configure methods for the opt system
func (s SuppressColumn1Flag) Configure(flags *Flags) { flags.SuppressColumn1 = s }
func (s SuppressColumn2Flag) Configure(flags *Flags) { flags.SuppressColumn2 = s }
func (s SuppressColumn3Flag) Configure(flags *Flags) { flags.SuppressColumn3 = s }
func (c CheckOrderFlag) Configure(flags *Flags)      { flags.CheckOrder = c }
func (t TotalFlag) Configure(flags *Flags)           { flags.Total = t }
