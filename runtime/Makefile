# Makefile for Koda Runtime Library
# Use a POSIX shell so `mkdir -p` / `rm -rf` work when GNU Make runs on Windows (cmd.exe does not).
ifeq ($(OS),Windows_NT)
SHELL := sh.exe
else
SHELL := sh
endif
.SHELLFLAGS := -ec

CC = gcc
CFLAGS = -O3 -fPIC -std=c11
AR = ar
ARFLAGS = rcs

SRCDIR = src
OBJDIR = obj
# Output: runtime/libkoda_runtime.a (linked by nativebuild)
LIBDIR = .

SRCS = $(SRCDIR)/value.c $(SRCDIR)/object.c $(SRCDIR)/gc.c $(SRCDIR)/koda_runtime.c
OBJS = $(patsubst $(SRCDIR)/%.c,$(OBJDIR)/%.o,$(SRCS))
TARGET = $(LIBDIR)/libkoda_runtime.a

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(OBJS)
	@mkdir -p $(LIBDIR)
	$(AR) $(ARFLAGS) $@ $^

$(OBJDIR)/%.o: $(SRCDIR)/%.c
	@mkdir -p $(OBJDIR)
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -rf $(OBJDIR) $(TARGET)
