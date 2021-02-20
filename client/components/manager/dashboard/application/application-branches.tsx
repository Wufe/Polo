import { IApplication, IApplicationBranchModel } from '@/state/models';
import { values } from 'mobx';
import { observer } from 'mobx-react-lite';
import React, { useEffect, useState } from 'react';
import dayjs from 'dayjs';

type TProps = {
    branches: IApplication['branches'];
    onSessionCreationSubmission: (checkout: string) => void;
}

export const ApplicationBranches = observer((props: TProps) => {

    const [openBranches, setOpenBranches] = useState<{[k:string]: boolean}>({});

    const toggleBranch = (name: string) => setOpenBranches(b => ({...b, [name]: !b[name]}));

    return <>
        <h4 className="my-2 text-gray-500 text-xs uppercase px-2 lg:px-6">Branches:</h4>
        <div className="divide-y dark:divide-gray-700">
            {sortBranches(Array.from(props.branches.values())).map((branch, key) =>
                <div
                    className="flex flex-col" key={key}>
                    <div className="flex items-end lg:items-center pt-1 pb-2 px-2 lg:px-6 cursor-pointer dark:hover:bg-nord-1" onClick={() => toggleBranch(branch.name)}>
                        <div className="flex-1">
                            <div className="text-xs text-gray-500 uppercase">Name</div>
                            <div className="text-sm whitespace-nowrap overflow-hidden overflow-ellipsis" title={branch.name}>{branch.name}</div>
                        </div>
                        <span className="text-center col-span-2">
                            <span className="text-sm underline cursor-pointer inline-block mx-3 hover:text-blue-400" onClick={() => props.onSessionCreationSubmission(branch.name)}>Create session</span>
                        </span>
                    </div>
                    <div className={`col-span-12 ${!openBranches[branch.name] && 'hidden'}`}>
                        <div className="grid grid-cols-12 gap-2 items-center p-3 lg:px-6 bg-nord5 dark:bg-nord1 rounded-sm">
                            <div className="col-span-12">
                                <div className="text-xs text-gray-500 uppercase">Message</div>
                                <div className="text-sm whitespace-nowrap overflow-hidden overflow-ellipsis">{branch.message}</div>
                            </div>
                            <div className="col-span-5">
                                <div className="text-xs text-gray-500 uppercase">Author</div>
                                <div className="text-sm whitespace-nowrap overflow-hidden overflow-ellipsis">{branch.author}</div>
                            </div>
                            <div className="col-span-7">
                                <div className="text-xs text-gray-500 uppercase">Date</div>
                                <div className="text-sm whitespace-nowrap overflow-hidden overflow-ellipsis">{dayjs(branch.date).format('DD MMM HH:mm')}</div>
                            </div>
                        </div>
                        
                    </div>
                    
                </div>)}
        </div>
        
    </>;
});

function sortBranches(branches: IApplicationBranchModel[]): IApplicationBranchModel[] {
    const preferred = [
        'master',
        'main',
        'hotfix',
        'develop',
        'dev',
        'feature'
    ];

    let result: IApplicationBranchModel[] = [];
    let length = branches.length;

    for (const pref of preferred) {
        for (let i = 0; i < length; i++) {
            const branch = branches[i];
            if (branch.name.toLowerCase().startsWith(pref.toLowerCase())) {
                result.push(branch);
                branches.splice(i, 1)
                length--;
                i--;
            }
        }
    }

    return result.concat(branches);
}