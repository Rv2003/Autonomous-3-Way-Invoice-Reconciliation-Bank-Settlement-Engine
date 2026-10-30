import { useCallback, useEffect, useState } from 'react';
import { fetchDocuments } from '../api/documents';

export default function useDocuments() {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [refreshVersion, setRefreshVersion] = useState(0);

  useEffect(() => {
    let active = true;

    fetchDocuments()
      .then((result) => {
        if (active) setDocuments(result);
      })
      .catch((requestError) => {
        if (active) {
          setError(requestError.message || 'Could not load documents from the backend.');
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [refreshVersion]);

  const loadDocuments = useCallback(() => {
    setLoading(true);
    setError('');
    setRefreshVersion((version) => version + 1);
  }, []);

  return { documents, loading, error, loadDocuments };
}
