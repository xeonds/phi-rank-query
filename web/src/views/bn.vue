<template>
  <div>
    <p v-if="is_pending">咕咕咕...查询中，请勿离开页面哦</p>
    <div class="title">
      <div class="l">
        <img src="/assets/Phigros_Icon_3.0.0.png" alt="icon">
        <div class="doc">
          <p>Phigros</p>
          <p>RankingScore查询</p>
        </div>
      </div>
      <div class="r">
        <p>玩家: {{ PlayerId }}</p>
        <p>Rks: {{ Rks }}</p>
        <div class="Challenge">
          <p>课题模式:</p>
          <div class="Challenge-r">
            <img :src="'/assets/' + ChallengeMode + '.png'" alt="Challenge">
            <p>{{ ChallengeModeRank }}</p>
          </div>
        </div>
        <p v-if="data">Data: {{ data }}</p>
        <p>同步时间: {{ lastDate }}</p>
      </div>
    </div>
    <div class="b19">
      <div class="L">
        <!-- phi -->
        <template v-for="(song, index) in b19_list" :key="song.num">
          <SongItem :index="'#' + (index + 1).toString()" :song="song" v-if="index % 2 === 0" />
        </template>

      </div>
      <div class="R">
        <template v-for="(song, index) in b19_list" :key="song.num">
          <SongItem :index="'#' + (index + 1).toString()" :song="song" v-if="index % 2 === 1" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Ref, onMounted, ref } from 'vue';
import SongItem from '@/components/song.vue';
import { Song } from '@/common';
import { getRating } from '@/common';

const PlayerId = ref('');
const Rks = ref('');
const ChallengeMode = ref('');
const ChallengeModeRank = ref('');
const data = ref('');
const lastDate = ref('');
const b19_list: Ref<Song[]> = ref([] as Song[]);
const sessionToken = ref('');

const askForSessionToken = () => {
  var token = localStorage.getItem('selectedSession');
  if (!token) {
    return false;
  } else {
    sessionToken.value = token!;
    return true;
  }
}
const fetchData = async (sessionToken: string) => {
  const cacheKey = 'history-' + sessionToken;
  // const cacheTimeout = 60000; // 1 minute in milliseconds
  const cacheTimeout = -1; // 1 minute in milliseconds

  const currentTime = Date.now();
  const cachedData = JSON.parse(localStorage.getItem(cacheKey) || 'null');

  if (cachedData && currentTime - cachedData.timestamp < cacheTimeout) {
    const latestData = cachedData.data[cachedData.data.length - 1];
    return { data: ref(latestData), err: ref(null) };
  } else {
    const data: Ref<any> = ref(null);
    const err: Ref<any> = ref(null);
    await fetch('/api/v1/bn', {
      method: 'POST',
      body: JSON.stringify({ session: sessionToken }),
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => response.json())
      .then(json => {
        data.value = json;
        const newData = {
          timestamp: currentTime,
          data: [...(cachedData?.data || []), json].slice(-32) // Max length of 32
        };
        localStorage.setItem(cacheKey, JSON.stringify(newData));
      })
      .catch(error => err.value = error);

    return { data, err };
  }
}
const parseData = (data: any) => {
  const getHiBit = (num: number) => {
    while (num > 10) {
      num = num / 10;
    }
    return num;
  }
  const getSubBit = (num: number) => {
    if (num >= 100) return num % 100;
    else if (num >= 10) return num % 10;
    else return num;
  }
  PlayerId.value = data.player || '';
  ChallengeMode.value = getHiBit(data.challenge).toFixed(0).toString() || '';
  ChallengeModeRank.value = getSubBit(data.challenge).toString() || '';
  lastDate.value = data.date || '';
  lastDate.value = new Date(lastDate.value).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).replace(/\//g, '-').replace(',', '');
  Rks.value = data.rks.toFixed(4).toString();
  b19_list.value = 
  [...data.phi, ...data.b27].map((item: any) => ({
    num: 0,
    song: item.Song || '',
    illustration: item.Illustration || '',
    rank: item.Level || '',
    difficulty: item.Difficulty || '',
    rks: item.Rks.toString() || '',
    Rating: getRating(item.FullCombo, item.Score) || '',
    score: item.Score.toString() || '',
    acc: (item.Acc + .005).toFixed(2).toString() || '',
    suggest: ''
  }));
}

const is_pending = ref(true);

onMounted(async () => {
  if (askForSessionToken()) {
    const { data, err } = await fetchData(sessionToken.value)
    if (err.value != null) {
      alert('查询失败，请重试');
      is_pending.value=false;
      window.location.href = '/#/';
    }
    parseData(data.value);
    is_pending.value=false;
  } else {
    alert('请选择Session');
    is_pending.value=false;
    window.location.href = '/#/session';
  }
});
</script>

<style lang="less" scoped>
.title {
  width: 100%;
  padding-top: 2rem;
  height: 12rem;
  display: flex;
  z-index: 1;
  flex-direction: row;
  padding-bottom: 2rem;
  align-items: center;
  justify-content: space-between;

  .l {
    background-color: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(3px);
    display: flex;
    flex-flow: row;
    width: 56%;
    height: 9.6rem;
    z-index: 2;
    writing-mode: horizontal-tb;
    clip-path: polygon(100% 0%, 90% 100%, 0% 100%, 0px 0px);
    display: flex;
    flex-direction: row;
    margin-right: -2rem;
    align-items: center;

    img {
      height: 75%;
      margin-inline: 1rem;
      border-radius: 2rem;
    }

    .doc {
      font-size: 1.4rem;
      margin: 10px;
    }
  }

  .r {
    background-color: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(3px);
    clip-path: polygon(100% 0%, 100% 100%, 0% 100%, 10% 0px);
    height: 8rem;
    width: 43.3%;
    display: flex;
    flex-direction: column;
    writing-mode: horizontal-tb;
    justify-content: center;
    margin-left: -4rem;
    padding-left: 3rem;
    z-index: 0;

    .Challenge {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: start;

      .Challenge-r {
        margin-inline-start: .5rem;
        padding: .1rem;
        display: flex;
        flex-flow: column;
        align-items: center;
        justify-content: center;

        img {
          width: 2.5rem;
          height: 2rem;
          margin-bottom: -1.8rem;
        }
      }
    }
  }
}

.b19 {
  display: flex;
  flex-flow: row;
  justify-content: space-evenly;

  .L {
    display: flex;
    flex-flow: column;
    align-items: center;
    padding-left: 2rem;
  }

  .R {
    display: flex;
    flex-flow: column;
    align-items: center;
    margin-top: 8%;
  }
}

.Nosignal,
.song {
  margin-bottom: 1rem;
}

.Nosignal {
  width: 87%;
  height: 243px;
  display: flex;
  position: relative;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  background-color: rgba(112, 112, 112, 0.4);

  .border_corner {
    z-index: 2500;
    position: absolute;
    width: 20px;
    height: 20px;
    background: rgba(0, 0, 0, 0);
    border: 6px solid #ffffff;

    &.border_corner_left_top {
      top: 5px;
      left: 5px;
      border-right: none;
      border-bottom: none;
      border-top-left-radius: 4px;
    }

    &.border_corner_right_top {
      top: 5px;
      right: 5px;
      border-left: none;
      border-bottom: none;
      border-top-right-radius: 4px;
    }

    &.border_corner_left_bottom {
      bottom: 5px;
      left: 5px;
      border-right: none;
      border-top: none;
      border-bottom-left-radius: 4px;
    }

    &.border_corner_right_bottom {
      bottom: 5px;
      right: 5px;
      border-left: none;
      border-top: none;
      border-bottom-right-radius: 4px;
    }
  }

  .line,
  .sqrt {
    filter: drop-shadow(0 0 5px #fff);
  }

  .timeout,
  .client {
    width: 90%;
    margin-top: -1%;
    margin-bottom: -1%;
    z-index: 5;
  }

  .timeout p {
    text-shadow:
      0 0 1px #d8f9ffab,
      1px 0 2px #d8f9ffab,
      -1px 0 2px #d8f9ffab,
      1px 0 5px #d8f9ffab,
      -1px 0 10px #d8f9ffab;
    margin-left: 5%;
    font-size: 75px;
    color: white;
  }

  .client p {
    text-shadow: 0 0 1px #8eeeff88, 1px 0 2px #8eeeff88, -1px 0 5px #8eeeff88;
    margin-left: 5%;
    font-size: 25px;
    color: #00B0F0;
  }

  .line {
    width: 20%;
    margin: 0;
    border-bottom: 5px solid;
    border-color: white;
    margin-left: 5%;
  }

  .sqrt {
    width: 89%;
    height: 20%;
    border: 3px solid;
    border-color: white;
    margin-left: 5%;
    margin-top: 16;
    z-index: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    background: repeating-linear-gradient(-45deg, #ffffff, #ffffff 5px, rgb(255 255 255 / 0%) 0, rgb(255 255 255 / 0%) 10px);
    filter: drop-shadow(0px 0px 3px #ffffff);

    p {
      color: white;
      font-size: 40px;
      filter: blur(3px);
    }
  }
}
</style>