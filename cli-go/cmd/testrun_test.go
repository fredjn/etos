package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/client"
	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ... existing code ... 