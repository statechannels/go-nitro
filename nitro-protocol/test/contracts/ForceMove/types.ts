export interface transitionType {
  whoSignedWhat: number[],
  appDatas: number[],
}

export interface testParams {
  largestTurnNum: number,
  support: transitionType,
  finalizesAt: number | undefined,
  reason: string | undefined
}
