import { check, fail, group } from 'k6';
import http from 'k6/http';

export const options = {
    vus: 1,
    thresholds: {
        // Ensure we have 100% compliance on API tests
        checks: [{ threshold: 'rate == 1.0', abortOnFail: true }],
    },
};

var targetProtocol = "http"
if (__ENV.PROTOCOL !== undefined) {
    targetProtocol = __ENV.PROTOCOL
}
var targetHost = "localhost"
if (__ENV.HOST !== undefined) {
    targetHost = __ENV.HOST
}
var targetPort = "80"
if (__ENV.PORT !== undefined) {
    targetPort = __ENV.PORT
}
const BASE_URL = `${targetProtocol}://${targetHost}:${targetPort}`;

export default () => {
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let testId = -1;
    const testName = `k6-${Date.now()}`;
    const testDesc = 'API Compliance Test';
    const testLat = 35.4183;
    const testLon = 76.5517;

    group('Initial listing check', function () {
        const placesRes = http.get(`${BASE_URL}/api/places`)
        check(placesRes, {
            'fetch returns appropriate status': (resp) => resp.status === 200,
        });

        // Confirm we do not have a place having the testName
        let places = placesRes.json();
        for (var i = 0; i < places.length; i++) {
            if (places[i].name === testName) {
                fail(`Test named "${testName}" already exists`);
            }
        }
    });

    group('Create a new place', function () {
        const createRes = http.post(`${BASE_URL}/api/places`, JSON.stringify({
            name: testName,
            description: testDesc,
            latitude: testLat,
            longitude: testLon,
        }), params);
        check(createRes, {
            'create returns appropriate status': (resp) => resp.status === 200,
            'and successfully creates a new place': (resp) => resp.json('id') !== '',
        });
        testId = createRes.json('id');
    });

    group('Retrieving a place', function () {
        const placeRes = http.get(`${BASE_URL}/api/places/${testId}`);
        check(placeRes, {
            'retrieving by id is successful': (resp) => resp.status === 200,
        });
        check(placeRes.json(), {
            'response provides attribute `id`': (place) => place.id === testId,
            'response provides attribute `name`': (place) => place.name === testName,
            'response provides attribute `description`': (place) => place.description === testDesc,
            'response provides attribute `latitude`': (place) => place.latitude === testLat,
            'response provides attribute `longitude`': (place) => place.longitude === testLon,
            'response provides attribute `created_at``': (place) => place.created_at !== undefined && place.created_at !== '',
            'response provides attribute `updated_at`': (place) => place.updated_at !== undefined && place.updated_at !== '',
        });
        // console.log("POST CREATE");
        // console.log(JSON.stringify(placeRes.body));

        // Ensure the place is returned in the list
        const placesRes = http.get(`${BASE_URL}/api/places`)
        let places = placesRes.json();
        let found = false;
        for (var i = 0; i < places.length; i++) {
            if (places[i].id === testId) {
                found = true;
                break;
            }
        }
        if (!found) {
            fail('Test place was not returned when retrieving all places');
        }
    });

    group('Update place by id', function () {
        const patchRes = http.patch(`${BASE_URL}/api/places/${testId}`, JSON.stringify({
            description: testDesc + " Updated"
        }), params);
        check(patchRes, {
            'update returns appropriate status': (resp) => resp.status === 200,
        });
        check(patchRes.json(), {
            'response provides attribute `id`': (place) => place.id === testId,
            'response provides attribute `name`': (place) => place.name === testName,
            'response provides modified attribute `description`': (place) => place.description === testDesc + " Updated",
            'response provides attribute `latitude`': (place) => place.latitude === testLat,
            'response provides attribute `longitude`': (place) => place.longitude === testLon,
            'response provides attribute `created_at``': (place) => place.created_at !== undefined && place.created_at !== '',
            'response provides attribute `updated_at`': (place) => place.updated_at !== undefined && place.updated_at !== '',
            'update changes modification date': (place) => place.updated_at !== place.created_at,
        });
        // console.log("POST UPDATE");
        // console.log(JSON.stringify(patchRes.body));
    });
    
    group('Delete place by id', function () {
        const deleteRes = http.del(`${BASE_URL}/api/places/${testId}`)
        check(deleteRes, {
            'delete returns appropriate status': (resp) => resp.status === 200,
        });
        // Confirm that the place has been removed
        const placeRes = http.get(`${BASE_URL}/api/places/${testId}`)
        check(placeRes, {
            'deleted place no longer available': (resp) => resp.status === 404,
        });
    });

}
