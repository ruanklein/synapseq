CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -g -D_POSIX_C_SOURCE=200809L
LDFLAGS = -lm

SRC_DIR = src
TEST_DIR = test
BUILD_DIR = output

UTIL_SOURCES = $(SRC_DIR)/utils/memory_utils.c $(SRC_DIR)/utils/string_utils.c
IO_SOURCES = $(SRC_DIR)/io/audio_output.c $(SRC_DIR)/io/wav_header.c
CORE_SOURCES = $(SRC_DIR)/core/time_manager.c $(SRC_DIR)/core/math_utils.c $(SRC_DIR)/core/wave_generator.c $(SRC_DIR)/core/audio_processor.c
PARSING_SOURCES = $(SRC_DIR)/parsing/sequence_parser.c
ALL_SOURCES = $(UTIL_SOURCES) $(IO_SOURCES) $(CORE_SOURCES) $(PARSING_SOURCES)

UTIL_OBJECTS = $(UTIL_SOURCES:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
IO_OBJECTS = $(IO_SOURCES:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
CORE_OBJECTS = $(CORE_SOURCES:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
PARSING_OBJECTS = $(PARSING_SOURCES:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
ALL_OBJECTS = $(ALL_SOURCES:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)

.PHONY: all clean test test-utils test-io test-time test-math test-parser test-wave

all: $(ALL_OBJECTS)

clean:
	rm -rf $(BUILD_DIR)
	rm -rf test_*

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/utils
	mkdir -p $(BUILD_DIR)/io
	mkdir -p $(BUILD_DIR)/core
	mkdir -p $(BUILD_DIR)/parsing

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c | $(BUILD_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

test: test-utils test-io test-time test-math test-parser test-wave test-audio

test-utils: $(UTIL_OBJECTS)
	@echo "=== Running Utils Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(TEST_DIR)/test_utils.c -o test_utils $(LDFLAGS)
	./test_utils

test-io: $(UTIL_OBJECTS) $(IO_OBJECTS)
	@echo "=== Running Audio Output Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(IO_OBJECTS) $(TEST_DIR)/test_audio_output.c -o test_audio_output $(LDFLAGS)
	./test_audio_output

test-time: $(UTIL_OBJECTS) $(CORE_OBJECTS)
	@echo "=== Running Time Manager Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(word 1,$(CORE_OBJECTS)) $(TEST_DIR)/test_time_manager.c -o test_time_manager $(LDFLAGS)
	./test_time_manager

test-math: $(UTIL_OBJECTS) $(CORE_OBJECTS)
	@echo "=== Running Math Utils Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(word 1,$(CORE_OBJECTS)) $(word 2,$(CORE_OBJECTS)) $(TEST_DIR)/test_math_utils.c -o test_math_utils $(LDFLAGS)
	./test_math_utils

test-parser: $(UTIL_OBJECTS) $(CORE_OBJECTS) $(PARSING_OBJECTS)
	@echo "=== Running Sequence Parser Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(word 1,$(CORE_OBJECTS)) $(word 2,$(CORE_OBJECTS)) $(PARSING_OBJECTS) $(TEST_DIR)/test_sequence_parser.c -o test_sequence_parser $(LDFLAGS)
	./test_sequence_parser

test-wave: $(UTIL_OBJECTS) $(CORE_OBJECTS)
	@echo "=== Running Wave Generator Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(word 2,$(CORE_OBJECTS)) $(word 3,$(CORE_OBJECTS)) $(TEST_DIR)/test_wave_generator.c -o test_wave_generator $(LDFLAGS)
	./test_wave_generator

test-audio: $(UTIL_OBJECTS) $(IO_OBJECTS) $(CORE_OBJECTS)
	@echo "=== Running Audio Processor Tests ==="
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(IO_OBJECTS) $(word 2,$(CORE_OBJECTS)) $(word 3,$(CORE_OBJECTS)) $(word 4,$(CORE_OBJECTS)) $(TEST_DIR)/test_audio_processor.c -o test_audio_processor $(LDFLAGS)
	./test_audio_processor

test-e2e: $(UTIL_OBJECTS) $(IO_OBJECTS) $(CORE_OBJECTS) $(PARSING_OBJECTS)
	@echo "=== Running End-to-End Integration Tests ==="
	@mkdir -p test/e2e
	$(CC) $(CFLAGS) $(UTIL_OBJECTS) $(IO_OBJECTS) $(CORE_OBJECTS) $(PARSING_OBJECTS) test/e2e/test_e2e_sequence_to_wav.c -o test_e2e_integration $(LDFLAGS)
	./test_e2e_integration