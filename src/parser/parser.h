// parser/parser.h
#ifndef SYNAPSEQ_PARSER_H
#define SYNAPSEQ_PARSER_H

#include <stdio.h>
#include "../core/config.h"
#include "../core/globals.h"

// Lexer functions (tokenization)
int readLine();
char *getWord();

// Main sequence parser
void readSeq(int ac, char **av);

// Specialized parsers
void readNameDef();
void readTimeLine();

// ReadNameDef utils
void normalizeAmplitude(Voice *voices, int numChannels, const char *line, int lineNum);
void checkBackgroundInSequence(NameDef *nd);
void handleOptionInSequence(char *p);

// Timeline utils
int voicesEq(Voice *v0, Voice *v1);
void correctPeriods();

// Error in sequence
void badSeq();
void badTime(char *tim);

int readTime(char *p, int *timp);

#endif // SYNAPSEQ_PARSER_H