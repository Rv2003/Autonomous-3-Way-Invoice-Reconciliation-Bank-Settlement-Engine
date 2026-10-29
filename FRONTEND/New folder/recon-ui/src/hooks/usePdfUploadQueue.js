import { useCallback, useRef, useState } from 'react';
import { displayDocumentStatus, uploadPdf } from '../api/documents';

const MAX_UPLOAD_SIZE = 25 * 1024 * 1024;
let nextQueueId = 0;

export default function usePdfUploadQueue() {
  const [queue, setQueue] = useState([]);
  const [processing, setProcessing] = useState(false);
  const queueRef = useRef([]);
  const processingRef = useRef(false);

  const updateQueue = useCallback((update) => {
    const next = update(queueRef.current);
    queueRef.current = next;
    setQueue(next);
  }, []);

  const addFiles = useCallback((fileList) => {
    const files = Array.from(fileList ?? []);
    if (!files.length) return;

    const entries = files.map((file) => {
      const error = !file.name.toLowerCase().endsWith('.pdf')
        ? 'Only PDF files can be uploaded.'
        : file.size > MAX_UPLOAD_SIZE
          ? 'PDF files must be 25 MB or smaller.'
          : '';
      return {
        id: ++nextQueueId,
        file: error ? null : file,
        name: file.name,
        status: error ? 'Failed' : 'Queued',
        message: error,
      };
    });

    updateQueue((current) => [...current, ...entries]);
  }, [updateQueue]);

  const processQueue = useCallback(async () => {
    if (processingRef.current) return;
    processingRef.current = true;
    setProcessing(true);

    try {
      while (true) {
        const next = queueRef.current.find((entry) => entry.status === 'Queued' && entry.file);
        if (!next) break;

        updateQueue((current) => current.map((entry) => (
          entry.id === next.id ? { ...entry, status: 'Uploading', message: '' } : entry
        )));

        try {
          const result = await uploadPdf(next.file);
          updateQueue((current) => current.map((entry) => (
            entry.id === next.id
              ? {
                ...entry,
                status: 'Uploaded',
                message: `Upload successful · Processing status: ${displayDocumentStatus(result.status)}.`,
              }
              : entry
          )));
        } catch (error) {
          updateQueue((current) => current.map((entry) => (
            entry.id === next.id
              ? { ...entry, status: 'Failed', message: error.message || 'Upload failed.' }
              : entry
          )));
        }
      }
    } finally {
      processingRef.current = false;
      setProcessing(false);
    }
  }, [updateQueue]);

  const retryFailed = useCallback(async () => {
    updateQueue((current) => current.map((entry) => (
      entry.status === 'Failed' && entry.file
        ? { ...entry, status: 'Queued', message: '' }
        : entry
    )));
    await processQueue();
  }, [processQueue, updateQueue]);

  const clearCompleted = useCallback(() => {
    if (processingRef.current) return;
    updateQueue((current) => current.filter((entry) => entry.status !== 'Uploaded'));
  }, [updateQueue]);

  const queuedCount = queue.filter((entry) => entry.status === 'Queued').length;
  const failedCount = queue.filter((entry) => entry.status === 'Failed').length;
  const retryableCount = queue.filter((entry) => entry.status === 'Failed' && entry.file).length;
  const completedCount = queue.filter((entry) => entry.status === 'Uploaded').length;

  const clearQueue = useCallback(() => {
    if (processingRef.current) return;
    updateQueue(() => []);
  }, [updateQueue]);

  return {
    queue,
    processing,
    queuedCount,
    failedCount,
    retryableCount,
    completedCount,
    addFiles,
    processQueue,
    retryFailed,
    clearCompleted,
    clearQueue,
  };
}
