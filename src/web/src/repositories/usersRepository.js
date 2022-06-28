import axios from "axios";
import resolveBaseURL from './resolveBaseURL'

const baseURL = resolveBaseURL(
    process.env.VUE_APP_USERS_SERVICE_DOMAIN,
    process.env.VUE_APP_USERS_SERVICE_PORT,
    process.env.VUE_APP_USERS_SERVICE_PATH
)

const connection = axios.create({
    baseURL
})

const resource = "/users";
export default {
    get(offset, count) {
        if (!offset) {
            offset = 0
        }
        if (!count) {
            count = 50
        }
        return connection.get(`${resource}/all?offset=${offset}&count=${count}`)
    },
    login(username, password) {
        let config = {
            auth: {
                username: username,
                password: password
            },
        }
        return connection.post(`/login`, {}, config)
    },
    getUnclaimedUser({ segment }) {
        return connection.get(`${resource}/unclaimed/?segment=${segment}`)
    },
    getRandomUser() {
        return connection.get(`${resource}/random/`)
    },
    getUserByID(userID) {
        if (!userID || userID.length == 0)
            throw "userID required"
        return connection.get(`${resource}/id/${userID}`)
    },
    getUserByUsername(username) {
        if (!username || username.length == 0)
            throw "username required"
        return connection.get(`${resource}/username/${username}`)
    },
    getUserByIdentityId(identityId) {
        if (!identityId || identityId.length == 0)
            throw "identityId required"
        return connection.get(`${resource}/identityid/${identityId}`)
    },
    createUser(user) {
        if (!user.username || user.username.length == 0)
            throw "username required"
        return connection.post(`${resource}`, user)
    },
    updateUser(user) {
        if (!user)
            throw "user required"
        return connection.put(`${resource}/id/${user.id}`, user)
    },
    claimUser(userId) {
        return connection.put(`${resource}/id/${userId}/claim`);
    },
    verifyAndUpdateUserPhoneNumber(userId, phoneNumber) {
        if (!userId || userId.length == 0)
            throw "userId required"
        let payload = {
            user_id: userId,
            phone_number: phoneNumber
        }
        return connection.put(`${resource}/id/${userId}/verifyphone`, payload)
    }
}