const segmenter = new Intl.Segmenter(['de', 'en'], { granularity: 'word' });

const collator = new Intl.Collator(['de', 'en'], { sensitivity: 'base' });

export function filterArray(array, searchTerm, fields) {
  searchTerm = searchTerm?.trim();
  if (!searchTerm) return array;

  return array.filter((item) => {
    const textToSearch = Object.keys(item)
      .map((key) => item[key])
      .join(' ');
    const segments = Array.from(segmenter.segment(textToSearch));

    return segments.some((s) => {
      if (!s.isWordLike) return false;

      if (collator.compare(s.segment, searchTerm) === 0) {
        return true;
      }

      const minLength = Math.min(s.segment.length, searchTerm.length);
      const subSegment = s.segment.substring(0, minLength);
      return collator.compare(subSegment, searchTerm) === 0;
    });
  });
}
