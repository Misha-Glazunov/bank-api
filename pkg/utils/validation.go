package utils

import (
    "regexp"
    "unicode"
)

func IsValidEmail(email string) bool {
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    return regexp.MustCompile(emailRegex).MatchString(email)
}

func IsStrongPassword(password string) bool {
    var (
        hasUpper   bool
        hasLower   bool
        hasNumber  bool
        hasSpecial bool
    )

    for _, c := range password {
        switch {
        case unicode.IsUpper(c):
            hasUpper = true
        case unicode.IsLower(c):
            hasLower = true
        case unicode.IsNumber(c):
            hasNumber = true
        case unicode.IsPunct(c) || unicode.IsSymbol(c):
            hasSpecial = true
        }
    }

    return hasUpper && hasLower && hasNumber && hasSpecial && len(password) >= 8
}

