#include "memory_utils.h"

#include <stdlib.h>
#include <string.h>
#include <stdio.h>

void* memory_alloc(size_t size) {
    void* ptr = calloc(1, size);
    if (!ptr) {
        fprintf(stderr, "Error: Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    return ptr;
}

void* memory_alloc_array(size_t count, size_t element_size) {
    return memory_alloc(count * element_size);
}

char* memory_duplicate_string(const char* str) {
    if (!str) return NULL;
    
    char* duplicate = strdup(str);
    if (!duplicate) {
        fprintf(stderr, "Error: Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    return duplicate;
}

void memory_free(void* ptr) {
    if (ptr) {
        free(ptr);
    }
}