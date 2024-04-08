<template>
  <div class="background">
    <div class="title">
      <div class="l">
        <img src="/assets/Phigros_Icon_3.0.0.png" alt="icon">
        <div class="doc">
          <p>Phigros</p>
          <p>B19查询</p>
        </div>
      </div>
      <div class="r">
        <p>Player: {{ PlayerId }}</p>
        <p>RankingScore: {{ Rks }}</p>
        <div class="Challenge">
          <p>ChallengeMode:</p>
          <div class="Challenge-r">
            <img :src="'/assets/' + ChallengeMode + '.png'" alt="Challenge">
            <p>{{ ChallengeModeRank }}</p>
          </div>
        </div>
        <p v-if="data">Data: {{ data }}</p>
        <p>Date: {{ Date }}</p>
      </div>
    </div>
    <div class="b19">
      <div class="L">
        <!-- phi -->
        <div v-if="phi.song" class="song">
          <div class="ill-box">
            <div class="num">
              <p>Phi</p>
            </div>
            <div class="ill">
              <img :src="phi.illustration" alt="ill">
            </div>
            <div :class="`rank-${phi.rank}`">
              <div class="org">
                <p>{{ phi.rank }}&ensp;{{ phi.difficulty }}</p>
              </div>
              <div class="rel">
                <p>{{ phi.rks }}</p>
              </div>
            </div>
          </div>
          <div class="info">
            <div class="songname">
              <p name="pvis">{{ phi.song }}</p>
            </div>
            <div class="songinfo">
              <div class="Rating">
                <img :src="'/assets/' + phi.Rating + '.png'" alt="Rating">
              </div>
              <div class="chengji">
                <div class="score">
                  <p>{{ phi.score }}</p>
                </div>
                <div class="line"></div>
                <div class="acc-box">
                  <div class="acc">
                    <p>{{ phi.acc }}%</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="Nosignal">
          <div class="border_corner border_corner_left_top"></div>
          <div class="border_corner border_corner_right_top"></div>
          <div class="border_corner border_corner_left_bottom"></div>
          <div class="border_corner border_corner_right_bottom"></div>
          <div class="line"></div>
          <div class="timeout">
            <p>TIME_OUT</p>
          </div>
          <div class="client">
            <p>>>> PhigrOS Client Finding Phi.score</p>
          </div>
          <div class="sqrt"></div>
        </div>

        <template v-for="(song, index) in b19_list" :key="song.num">
          <div class="song" v-if="index % 2 === 1">
            <div class="ill-box">
              <div class="num">
                <p name="pvis">#{{ index + 1 }}</p>
              </div>
              <div class="ill">
                <img :src="song.illustration" alt="ill">
              </div>
              <div :class="`rank-${song.rank}`">
                <div class="org">
                  <p>{{ song.rank }}&ensp;{{ song.difficulty }}</p>
                </div>
                <div class="rel">
                  <p>{{ song.rks }}</p>
                </div>
              </div>
            </div>
            <div class="info">
              <div class="songname">
                <p name="pvis">{{ song.song }}</p>
              </div>
              <div class="songinfo">
                <div class="Rating">
                  <img :src="'/assets/' + song.Rating + '.png'" alt="Rating">
                </div>
                <div class="chengji">
                  <div class="score">
                    <p>{{ song.score }}</p>
                  </div>
                  <div class="line"></div>
                  <div class="acc-box">
                    <div class="acc">
                      <p>{{ song.acc }}%</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>

      </div>
      <div class="R">
        <template v-for="(song, index) in b19_list" :key="song.num">
          <div class="song" v-if="index % 2 === 0">
            <div class="ill-box">
              <div class="num">
                <p name="pvis">#{{ index + 1 }}</p>
              </div>
              <div class="ill">
                <img :src="song.illustration" alt="ill">
              </div>
              <div :class="`rank-${song.rank}`">
                <div class="org">
                  <p>{{ song.rank }}&ensp;{{ song.difficulty }}</p>
                </div>
                <div class="rel">
                  <p>{{ song.rks }}</p>
                </div>
              </div>
            </div>
            <div class="info">
              <div class="songname">
                <p name="pvis">{{ song.song }}</p>
              </div>
              <div class="songinfo">
                <div class="Rating">
                  <img :src="'/assets/' + song.Rating + '.png'" alt="Rating">
                </div>
                <div class="chengji">
                  <div class="score">
                    <p>{{ song.score }}</p>
                  </div>
                  <div class="line"></div>
                  <div class="acc-box">
                    <div class="acc">
                      <p>{{ song.acc }}%</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Ref, onMounted, ref } from 'vue';

interface Song {
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

const PlayerId = ref('');
const Rks = ref('');
const ChallengeMode = ref('');
const ChallengeModeRank = ref('');
const data = ref('');
const Date = ref('');
const phi = ref({
  song: '',
  illustration: '',
  rank: '',
  difficulty: '',
  rks: '',
  Rating: '',
  score: '',
  acc: '',
  suggest: ''
});
const b19_list: Ref<Song[]> = ref([] as Song[]);

const askForSessionToken = () => {
  var sessionToken = prompt("Please enter your session token:");
  if (sessionToken) {
    fetch('/api/v1/b19', {
      method: 'POST',
      body: JSON.stringify({ session: sessionToken }),
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => response.json())
      .then(data => {
        parseData(data);
      })
      .catch(error => {
        console.error('Error:', error);
      });
  }
}
const getRating = (score: number) => {
  if (score >= 1000000) return 'phi';
  if (score >= 960000) return 'V';
  if (score >= 920000) return 'S';
  if (score >= 880000) return 'A';
  if (score >= 820000) return 'B';
  if (score >= 700000) return 'C';
  return 'F';
}
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
const parseData = (data: any) => {
  PlayerId.value = data.player || '';
  ChallengeMode.value = getHiBit(data.challenge).toFixed(0).toString() || '';
  ChallengeModeRank.value = getSubBit(data.challenge).toString() || '';
  Date.value = data.date || '';
  Rks.value = data.rks.toFixed(4).toString();
  phi.value.song = data.phi.Song || '';
  phi.value.illustration = data.phi.Illustration || '';
  phi.value.rank = data.phi.Level || '';
  phi.value.difficulty = data.phi.Difficulty || '';
  phi.value.rks = data.phi.Rks.toString() || '';
  phi.value.Rating = data.phi.Score >= 1000000 ? 'phi' : '';
  phi.value.score = data.phi.Score.toString() || '';
  phi.value.acc = data.phi.Acc.toString() || '';
  phi.value.suggest = '';

  b19_list.value = data.b19.map((item: any) => ({
    num: 0,
    song: item.Song || '',
    illustration: item.Illustration || '',
    rank: item.Level || '',
    difficulty: item.Difficulty || '',
    rks: item.Rks.toString() || '',
    Rating: getRating(item.Score) || '',
    score: item.Score.toString() || '',
    acc: item.Acc.toString() || '',
    suggest: ''
  }));
}

onMounted(() => {
  askForSessionToken();
});
</script>

<style lang="less" scoped>
body {
  margin: 0;
  background: url('/assets/Star1.png') center no-repeat;
  background-size: 100% 100%;

  p {
    color: white;
    margin: 0px;
    margin-block: 0;
    margin-inline: 0;
  }
}

.title {
  width: 100%;
  margin-top: 4rem;
  height: 12rem;
  display: flex;
  z-index: 1;
  flex-direction: row;
  margin-bottom: 100px;
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

.song {
  width: auto;
  height: 12rem;
  display: flex;
  flex-flow: row;
  align-items: center;
  justify-content: center;

  .ill-box {
    width: auto;
    height: inherit;
    display: flex;
    flex-flow: column;
    align-items: flex-start;
    justify-content: center;
    z-index: 1;
    filter: drop-shadow(2px 4px 10px #00000055);

    .num {
      p {
        font-size: 1rem;
        margin: 0;
        color: black;
      }

      z-index: 2;
      width: 2.4rem;
      height: 1.8rem;
      display: flex;
      background-color: rgba(255, 255, 255, 1);
      clip-path: polygon(100% 0, 80% 100%, 0 100%, 20% 0);
      justify-content: center;
      align-items: center;
    }

    .ill {
      width: 12.8rem;
      height: auto;
      display: flex;
      flex-direction: column;
      margin-top: -1.8rem;
      margin-bottom: -3.2rem;
      clip-path: polygon(100% 0, 90% 100%, 0 100%, 10% 0);
      overflow: hidden;
    }

    .rank-EZ {
      background-color: #92D050;
    }

    .rank-HD {
      background-color: #00B0F0;
    }

    .rank-IN {
      background-color: rgb(255, 0, 0);
    }

    .rank-AT {
      background-color: rgb(110, 110, 110);
    }

    .rank-EZ,
    .rank-HD,
    .rank-IN,
    .rank-AT {
      width: 4.8rem;
      height: 3.2rem;
      display: flex;
      align-items: center;
      flex-flow: column;
      justify-content: center;
      margin-left: -15%;
      clip-path: polygon(100% 0, 85% 100%, 0 100%, 15% 0);

      .org,
      .rel {
        display: flex;
        width: auto;
        justify-content: center;
        align-items: center;
      }

      .org {
        height: 35%;
        margin-left: 10%;
        font-size: 0.8rem;
      }

      .rel {
        height: 49%;
        margin-right: 10%;
        font-size: 1.2rem;
      }
    }
  }

  .info {
    margin-left: -3rem;
    padding-left: 3rem;
    padding-top: .2rem;
    padding-bottom: .3rem;
    display: flex;
    flex-flow: column;
    align-items: center;
    justify-content: center;
    width: 14rem;
    height: 6rem;
    background-color: rgba(0, 0, 0, 0.5);
    clip-path: polygon(100% 0, 90% 100%, 0 100%, 10% 0);

    .songname {
      width: 90%;
      height: 25%;
      display: flex;
      justify-content: center;
      align-items: center;

      p {
        font-size: 1rem;
        margin: 0;
        text-align: center;
      }
    }

    .songinfo {
      width: 100%;
      height: 70%;
      display: flex;
      flex-flow: row;
      align-items: center;

      .Rating {
        width: 3rem;
        height: 3rem;
        display: flex;

        img {
          width: 100%;
          filter: drop-shadow(0px 0px 10px #ffffffbb);
        }
      }

      .chengji {
        width: 80%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        .score {
          font-size: 1.6rem;
        }

        .line {
          width: 90%;
          margin: 0;
          border-bottom: 1px solid;
          border-color: white;
        }

        .acc-box {
          display: flex;
          justify-content: center;
          align-items: flex-start;

          .acc {
            display: flex;
            justify-content: flex-end;
            align-items: center;

            p {
              font-size: 1.2rem;
            }
          }

          .suggest {
            display: none;
            width: auto;
            display: flex;

            p {
              font-size: 1rem;
              color: rgb(192, 221, 173);
            }
          }
        }
      }
    }
  }
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