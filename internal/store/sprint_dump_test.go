package store_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kulti/task-list-bot/internal/store"
)

func TestStoreSprintDump(t *testing.T) {
	t.Parallel()

	testdataDir := filepath.Join("testdata", "sprint_dump")
	fileInfos, err := ioutil.ReadDir(testdataDir)
	require.NoError(t, err)

	for _, fi := range fileInfos { //nolint:paralleltest,lll // false-positive https://github.com/kunwardeep/paralleltest/issues/8
		fi := fi
		require.True(t, fi.IsDir())
		t.Run(fi.Name(), func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			t.Cleanup(mockCtrl.Finish)

			mockRepository := NewMockRepository(mockCtrl)
			store := store.New(mockRepository)

			sprintData, err := ioutil.ReadFile(filepath.Join(testdataDir, fi.Name(), "in"))
			require.NoError(t, err)
			mockRepository.EXPECT().CurrentSprint().Return(sprintData, nil)

			actual, err := store.CurrentSprintDump()
			require.NoError(t, err)

			expectedFileName := filepath.Join(testdataDir, fi.Name(), "out")
			if update {
				require.NoError(t, ioutil.WriteFile(expectedFileName, []byte(actual), 0600))
			}

			expectedBytes, err := ioutil.ReadFile(expectedFileName)
			require.NoError(t, err)

			require.Equal(t, string(expectedBytes), actual)
		})
	}
}
