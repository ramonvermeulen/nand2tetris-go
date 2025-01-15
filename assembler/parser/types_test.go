package parser

import "testing"

func TestAInstruction_String(t *testing.T) {
	aInstruction := AInstruction{Symbol: "123"}
	expected := "@123\n"
	if aInstruction.String() != expected {
		t.Errorf("Expected %s, got %s", expected, aInstruction.String())
	}
}

func TestCInstruction_String(t *testing.T) {
	// TODO(ramon) add more cases (parameterized tests)
	cInstruction := CInstruction{Dest: "D", Comp: "M", Jump: "JGT"}
	expected := "D=M;JGT\n"
	if cInstruction.String() != expected {
		t.Errorf("Expected %s, got %s", expected, cInstruction.String())
	}
}

func TestLabel_String(t *testing.T) {
	lInstruction := Label{Name: "LOOP"}
	expected := "(LOOP)\n"
	if lInstruction.String() != expected {
		t.Errorf("Expected %s, got %s", expected, lInstruction.String())
	}
}
