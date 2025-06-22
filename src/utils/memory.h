// utils/memory.h
#ifndef SYNAPSEQ_MEMORY_H
#define SYNAPSEQ_MEMORY_H

#include <stddef.h>

// Memory allocation functions
void *Alloc(size_t len);
char *StrDup(char *str);

#endif // SYNAPSEQ_MEMORY_H