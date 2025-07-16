#include "string_utils.h"

#include <string.h>
#include <ctype.h>
#include <stdlib.h>

char* string_trim(char* str) {
    if (!str) return NULL;
    
    while (isspace(*str)) str++;
    
    if (*str == 0) return str;
    
    char* end = str + strlen(str) - 1;
    while (end > str && isspace(*end)) end--;
    
    end[1] = '\0';
    
    return str;
}

bool string_is_empty(const char* str) {
    if (!str) return true;
    
    while (*str) {
        if (!isspace(*str)) return false;
        str++;
    }
    return true;
}

char* string_get_next_word(char** str) {
    if (!str || !*str) return NULL;
    
    while (isspace(**str)) (*str)++;
    
    if (!**str) return NULL;
    
    char* start = *str;
    
    while (**str && !isspace(**str)) (*str)++;
    
    char* end = *str;
    if (*end) {
        *end = '\0';
        (*str)++;
    }
    
    return start;
}

void string_to_lowercase(char* str) {
    if (!str) return;
    
    while (*str) {
        *str = tolower(*str);
        str++;
    }
}

bool string_is_valid_name(const char* str) {
    if (!str || !*str) return false;
    
    if (!isalpha(*str)) return false;
    
    str++;
    while (*str) {
        if (!isalnum(*str) && *str != '_' && *str != '-') {
            return false;
        }
        str++;
    }
    
    return true;
}