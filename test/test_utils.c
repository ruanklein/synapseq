#include "../src/utils/memory_utils.h"
#include "../src/utils/string_utils.h"

#include <stdio.h>
#include <assert.h>
#include <string.h>

void test_memory_utils() {
    printf("Testing memory_utils...\n");
    
    void* ptr = memory_alloc(100);
    assert(ptr != NULL);
    memory_free(ptr);
    
    char* str = memory_duplicate_string("hello");
    assert(strcmp(str, "hello") == 0);
    memory_free(str);
    
    printf("✓ memory_utils passed\n");
}

void test_string_utils() {
    printf("Testing string_utils...\n");
    
    char test_str[] = "  hello world  ";
    char* trimmed = string_trim(test_str);
    assert(strcmp(trimmed, "hello world") == 0);
    
    assert(string_is_empty("   ") == true);
    assert(string_is_empty("hello") == false);
    
    assert(string_is_valid_name("test_preset") == true);
    assert(string_is_valid_name("123invalid") == false);
    
    printf("✓ string_utils passed\n");
}

int main() {
    printf("=== Testing Basic Utilities ===\n");
    
    test_memory_utils();
    test_string_utils();
    
    printf("\n✓ All tests passed!\n");
    return 0;
}