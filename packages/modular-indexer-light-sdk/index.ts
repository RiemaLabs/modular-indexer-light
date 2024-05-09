import Worker from "./worker?worker&inline";
import {SDK} from "./worker";
import {wrap} from "comlink";

const LinkedSDK = wrap<typeof SDK>(new Worker());

export async function create(): Promise<SDK> {
    return await new LinkedSDK();
}

export type * from "./worker";
