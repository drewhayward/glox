package lox 

type Value interface {
}


func boolify(v any) bool {
    if v == nil || v == false {
        return false
    }

    return true
}


func Evaluate(node Node) Value {
    switch nt := node.(type) {

    case UnaryExpr:
        value := Evaluate(nt)
        if nt.Operation == BANG {
            // Assert that the value is a bool?
            return boolify(value)
        } else if nt.Operation == MINUS {
            
        }
    case GroupingExpr:

    }

    panic("Ahhhh")
}
