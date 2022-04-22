import {writeFileSync} from 'fs';

const jsonPath = __dirname + '/../addresses.json';
// eslint-disable-next-line @typescript-eslint/no-var-requires
const addresses = require(jsonPath);

function deepDelete(object: any, keyToDelete: string) {
  Object.keys(object).forEach(key => {
    if (key === keyToDelete) delete object[key];
    else if (typeof object[key] === 'object') deepDelete(object[key], keyToDelete);
  });
}
const keyToDelete = 'abi';
deepDelete(addresses, keyToDelete);
writeFileSync(jsonPath, JSON.stringify(addresses, null, 2));
