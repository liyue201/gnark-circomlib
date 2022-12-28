package circuits

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func Window4(api frontend.API, in []frontend.Variable, base []frontend.Variable) ([]frontend.Variable, []frontend.Variable) {
	if len(in) != 4 || len(base) != 2 {
		panic("invalid param")
	}

	out := make([]frontend.Variable, 2)
	out8 := make([]frontend.Variable, 2)

	c := Make2DVariableArray(2, 8)

	// in[0]  -> 1*BASE
	c[0][0] = base[0]
	c[1][0] = base[1]

	// in[1] -> 2*BASE
	dbl2 := MontgomeryDouble(api, base[:2])
	c[0][1] = dbl2[0]
	c[1][1] = dbl2[1]

	// in[2] -> 3*BASE
	adr3 := MontgomeryAdd(api, base[:2], dbl2[:2])
	c[0][2] = adr3[0]
	c[1][2] = adr3[1]

	// in[3] -> 4*BASE
	adr4 := MontgomeryAdd(api, base[:2], adr3[:2])
	c[0][3] = adr4[0]
	c[1][3] = adr4[1]

	// in[4] -> 5*BASE
	adr5 := MontgomeryAdd(api, base[:2], adr4[:2])
	c[0][4] = adr5[0]
	c[1][4] = adr5[1]

	// in[5] -> 6*BASE
	adr6 := MontgomeryAdd(api, base[:2], adr5[:2])
	c[0][5] = adr6[0]
	c[1][5] = adr6[1]

	// in[6] -> 7*BASE
	adr7 := MontgomeryAdd(api, base[:2], adr6[:2])
	c[0][6] = adr7[0]
	c[1][6] = adr7[1]

	// in[7] -> 8*BASE
	adr8 := MontgomeryAdd(api, base[:2], adr7[:2])
	c[0][7] = adr8[0]
	c[1][7] = adr8[1]

	out8 = adr8

	mux := MultiMux3(api, c, in[:3])
	out[0] = mux[0]
	out[1] = api.Sub(mux[1], api.Mul(mux[1], 2, in[3]))

	return out, out8
}

func Segment(api frontend.API, in []frontend.Variable, base []frontend.Variable) []frontend.Variable {
	if len(base) != 2 {
		panic("invalid param")
	}
	nWindows := len(in) / 4
	out := make([]frontend.Variable, 2)

	e2m := Edwards2Montgomery(api, base)

	windowsOut := MakeVariableArray(2)
	windowsOut8 := MakeVariableArray(2)
	adders := MakeVariableArray(2)
	for i := 0; i < nWindows; i++ {
		windowsIn := MakeVariableArray(4)
		for j := 0; j < 4; j++ {
			windowsIn[j] = in[i*4+j]
		}
		if i == 0 {
			windowsOut, windowsOut8 = Window4(api, windowsIn, e2m)
		} else {
			doublers1 := MontgomeryDouble(api, windowsOut8)
			doublers2 := MontgomeryDouble(api, doublers1)

			thisWindowsOut, thisWindowsOut8 := Window4(api, windowsIn, doublers2)

			if i == 1 {
				adders = MontgomeryAdd(api, windowsOut, thisWindowsOut)
			} else {
				adders = MontgomeryAdd(api, adders, thisWindowsOut)
			}
			windowsOut = thisWindowsOut
			windowsOut8 = thisWindowsOut8
		}
	}
	if nWindows > 1 {
		out = Montgomery2Edwards(api, adders)
	} else {
		out = Montgomery2Edwards(api, windowsOut)
	}
	return out
}

func Pedersen(api frontend.API, in []frontend.Variable) []frontend.Variable {

	var BASE [10][2]*big.Int
	BASE[0][0], _ = new(big.Int).SetString("10457101036533406547632367118273992217979173478358440826365724437999023779287", 10)
	BASE[0][1], _ = new(big.Int).SetString("19824078218392094440610104313265183977899662750282163392862422243483260492317", 10)
	BASE[1][0], _ = new(big.Int).SetString("2671756056509184035029146175565761955751135805354291559563293617232983272177", 10)
	BASE[1][1], _ = new(big.Int).SetString("2663205510731142763556352975002641716101654201788071096152948830924149045094", 10)
	BASE[2][0], _ = new(big.Int).SetString("5802099305472655231388284418920769829666717045250560929368476121199858275951", 10)
	BASE[2][1], _ = new(big.Int).SetString("5980429700218124965372158798884772646841287887664001482443826541541529227896", 10)
	BASE[3][0], _ = new(big.Int).SetString("7107336197374528537877327281242680114152313102022415488494307685842428166594", 10)
	BASE[3][1], _ = new(big.Int).SetString("2857869773864086953506483169737724679646433914307247183624878062391496185654", 10)
	BASE[4][0], _ = new(big.Int).SetString("20265828622013100949498132415626198973119240347465898028410217039057588424236", 10)
	BASE[4][1], _ = new(big.Int).SetString("1160461593266035632937973507065134938065359936056410650153315956301179689506", 10)
	BASE[5][0], _ = new(big.Int).SetString("1487999857809287756929114517587739322941449154962237464737694709326309567994", 10)
	BASE[5][1], _ = new(big.Int).SetString("14017256862867289575056460215526364897734808720610101650676790868051368668003", 10)
	BASE[6][0], _ = new(big.Int).SetString("14618644331049802168996997831720384953259095788558646464435263343433563860015", 10)
	BASE[6][1], _ = new(big.Int).SetString("13115243279999696210147231297848654998887864576952244320558158620692603342236", 10)
	BASE[7][0], _ = new(big.Int).SetString("6814338563135591367010655964669793483652536871717891893032616415581401894627", 10)
	BASE[7][1], _ = new(big.Int).SetString("13660303521961041205824633772157003587453809761793065294055279768121314853695", 10)
	BASE[8][0], _ = new(big.Int).SetString("3571615583211663069428808372184817973703476260057504149923239576077102575715", 10)
	BASE[8][1], _ = new(big.Int).SetString("11981351099832644138306422070127357074117642951423551606012551622164230222506", 10)
	BASE[9][0], _ = new(big.Int).SetString("18597552580465440374022635246985743886550544261632147935254624835147509493269", 10)
	BASE[9][1], _ = new(big.Int).SetString("6753322320275422086923032033899357299485124665258735666995435957890214041481", 10)

	n := len(in)
	nSegments := (n + 199) / 200

	segments := Make2DVariableArray(nSegments, 2)
	for i := 0; i < nSegments; i++ {
		nBits := 200
		if i == nSegments-1 {
			nBits = n - (nSegments-1)*200
		}
		nWindows := (nBits + 3) / 4

		segmentsIn := MakeVariableArray(nWindows * 4)
		for j := 0; j < nBits; j++ {
			segmentsIn[j] = in[i*200+j]
		}
		segments[i] = Segment(api, segmentsIn, []frontend.Variable{BASE[i][0], BASE[i][1]})
	}

	var adderX, adderY frontend.Variable
	for i := 0; i < nSegments-1; i++ {
		if i == 0 {
			adderX, adderY = BabyAdd(api, segments[0][0], segments[0][1], segments[1][0], segments[1][1])
		} else {
			adderX, adderY = BabyAdd(api, adderX, adderY, segments[i+1][0], segments[i+1][1])
		}
	}

	out := MakeVariableArray(2)
	if nSegments > 1 {
		out[0] = adderX
		out[1] = adderY
	} else {
		out[0] = segments[0][0]
		out[1] = segments[0][1]
	}
	return out
}
