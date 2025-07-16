#ifndef STRING_UTILS_H
#define STRING_UTILS_H

#include <stdbool.h>

char* string_trim(char* str);

bool string_is_empty(const char* str);

char* string_get_next_word(char** str);

void string_to_lowercase(char* str);

bool string_is_valid_name(const char* str);

#endif // STRING_UTILS_H