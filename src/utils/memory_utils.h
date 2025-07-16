#ifndef MEMORY_UTILS_H
#define MEMORY_UTILS_H

#include <stddef.h>

void* memory_alloc(size_t size);

void* memory_alloc_array(size_t count, size_t element_size);

char* memory_duplicate_string(const char* str);

void memory_free(void* ptr);

#endif // MEMORY_UTILS_H