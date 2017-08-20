package main

import "github.com/stretchr/testify/assert"
import "testing"

func TestFetchAndScanEmptyUrls(t *testing.T) {
	fetchAndScan := NewFetchAndScan(0, 0, 0)
	data, e := fetchAndScan.Do([]string{})
	assert.NoError(t, e)
	assert.Equal(t, []FetchAndScanData{}, data)
}

func makeFuncAndScanContext(
	fetchWorkers, scanWorkers, taskBufferSize int,
)(
	fetchAndScan *FetchAndScan,
	context      *FetchAndScanContext,
	done         chan *FetchAndScanTask,
){
	fetchAndScan = NewFetchAndScan(fetchWorkers, scanWorkers, taskBufferSize)
	done         = make(chan *FetchAndScanTask)
	context      = &FetchAndScanContext{done, 0}
	return
}

func TestFetchWorkerShouldAdvance(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(2, 0, 0)
	for i := 0; i < 3; i++ {
		fs.fetchPipe <- NewFetchAndScanTask(context, "")
		assert.Error(t, (<-done).e)
	}
}

func TestFetchWorkerShouldAdvanceCanceled(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(2, 0, 0)
	context.Cancel()
	for i := 0; i < 3; i++ {
		fs.fetchPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}

func TestScanWorkerShouldAdvance(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(0, 2, 0)
	for i := 0; i < 3; i++ {
		fs.scanPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}

func TestScanWorkerShouldAdvanceCanceled(t *testing.T) {
	fs, context, done := makeFuncAndScanContext(0, 2, 0)
	context.Cancel()
	for i := 0; i < 3; i++ {
		fs.scanPipe <- NewFetchAndScanTask(context, "")
		assert.NoError(t, (<-done).e)
	}
}
