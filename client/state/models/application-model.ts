import { APIPayload } from "@/api/common";
import { createNewSessionAPI } from "@/api/applications";
import { flow, Instance, SnapshotIn, SnapshotOut, types } from "mobx-state-tree";
import { ISession } from "./session-model";

export const ApplicationBranchModel = types.model({
    name   : types.string,
    hash   : types.string,
    author : types.string,
    date   : types.string,
    message: types.string,
})

export interface IApplicationBranchModel extends Instance<typeof ApplicationBranchModel> {}

export const ApplicationModel = types.model({
    name                 : types.string,
    remote               : types.string,
    target               : types.string,
    host                 : types.string,
    maxConcurrentSessions: types.number,
    folder               : types.string,
    branches             : types.map(ApplicationBranchModel)
})
.actions(self => {

    const newSession = flow(function* newSession(checkout: string) {
        const session: APIPayload<ISession> = yield  createNewSessionAPI(self.name, checkout);
        return session;
    });

    return { newSession };
})

export interface IApplication extends Instance<typeof ApplicationModel> { }
export interface IApplicationSnapshotOut extends SnapshotOut<typeof ApplicationModel> { }
export interface IApplicationSnapshotIn extends SnapshotIn<typeof ApplicationModel> { }