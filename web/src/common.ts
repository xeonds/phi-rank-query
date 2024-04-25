export interface Song {
  num: number;
  song: string;
  illustration: string;
  rank: string;
  difficulty: string;
  rks: string;
  Rating: string;
  score: string;
  acc: string;
  suggest: string;
}

export interface RankTableItem {
    title: string;
    difficulty: string;
    rank: number;
}

export interface RankTableItemShrink {
    title: string;
    difficulty: {
      EZ: number;
      HD: number;
      IN: number;
      AT?: number;
    }
}

export const getRating = (fc: boolean, score: number) => {
  if (score >= 1000000) return "phi";
  if (fc) return "FC";
  if (score >= 960000) return "V";
  if (score >= 920000) return "S";
  if (score >= 880000) return "A";
  if (score >= 820000) return "B";
  if (score >= 700000) return "C";
  return "F";
};

export const parseSongList = (songList: any) => {
  let result: Song[] = [] as Song[];
  result = songList.map((item: any) => ({
    num: 0,
    song: item.Song || "",
    illustration: item.Illustration || "",
    rank: item.Level || "",
    difficulty: item.Difficulty || "",
    rks: item.Rks.toString() || "",
    Rating: getRating(item.FullCombo, item.Score) || "",
    score: item.Score.toString() || "",
    acc: (item.Acc + 0.005).toFixed(2).toString() || "",
    suggest: "",
  }));
  return result;
};

export const parseSong = (song: any) => {
  return {
    num: 0,
    song: song.Song || "",
    illustration: song.Illustration || "",
    rank: song.Level || "",
    difficulty: song.Difficulty || "",
    rks: song.Rks.toString() || "",
    Rating: getRating(song.FullCombo, song.Score) || "",
    score: song.Score.toString() || "",
    acc: (song.Acc + 0.005).toFixed(2).toString() || "",
    suggest: "",
  };
};
