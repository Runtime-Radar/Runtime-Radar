type OneOnly<O, K extends keyof O> = { [key in Exclude<keyof O, K>]: null } & Pick<O, K>;

type OneOfByKey<O> = { [key in keyof O]: OneOnly<O, key> };

type ValueOf<O> = O[keyof O];

export type OneOf<O> = ValueOf<OneOfByKey<O>>;
