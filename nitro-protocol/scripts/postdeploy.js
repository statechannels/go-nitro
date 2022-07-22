// eslint-disable-next-line no-undef
const {writeFileSync} = require('fs');

// eslint-disable-next-line no-undef
const jsonPath = __dirname + '/../addresses.json';
// eslint-disable-next-line no-undef
const addresses = require(jsonPath);

function deepDelete(object, keyToDelete) {
  Object.keys(object).forEach(key => {
    if (key === keyToDelete) delete object[key];
    else if (typeof object[key] === 'object') deepDelete(object[key], keyToDelete);
  });
}
const keyToDelete = 'abi';
deepDelete(addresses, keyToDelete);
writeFileSync(jsonPath, JSON.stringify(addresses, null, 2));
