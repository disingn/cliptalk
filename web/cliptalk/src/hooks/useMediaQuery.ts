import { useState, useEffect } from 'react';

function useMediaQuery(query: string) {
    const [matches, setMatches] = useState(window.matchMedia(query).matches);

    useEffect(() => {
        const mediaQueryList = window.matchMedia(query);
        const documentChangeHandler = () => setMatches(mediaQueryList.matches);

        mediaQueryList.addListener(documentChangeHandler);
        documentChangeHandler();
        return () => {
            mediaQueryList.removeListener(documentChangeHandler);
        };
    }, [query]);

    return matches;
}

export default useMediaQuery;