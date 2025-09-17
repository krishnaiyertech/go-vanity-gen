# SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
# SPDX-License-Identifier: Apache-2.0

GO_LINT=golangci-lint

.PHONY: init

go.lint:
	${GO_LINT} run