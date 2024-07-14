function main() {
    let step = document.getElementById("step");

    let ip = new Register("ip");
    let r0 = new Register("r0");
    let r1 = new Register("r1");
    let r2 = new Register("r2");
    let r3 = new Register("r3");

    step.addEventListener("click", () => {
        vmStep().then((state) => {
            ip.changeValue(state.regs.ip);
            r0.changeValue(state.regs.r[0]);
            r1.changeValue(state.regs.r[1]);
            r2.changeValue(state.regs.r[2]);
            r3.changeValue(state.regs.r[3]);
        });
    });
}

/**
 * Send a request to make a single vm step.
 * 
 * @returns {Promise} VM state after the step.
 */
function vmStep() {
    return fetch("/vm/step", {method: "POST"})
        .then((response) => response.json());
}

/**
 * Formats a given byte number as string of hex digits. Values less than 0x10
 * are prefixed with 0, thus function always returns a string with 2 digits.
 * 
 * @param {number} b Byte value. Must be a number in range 0 <= b <= 255. 
 * @returns {string} Byte value in hex format. 
 */
function formatByte(b) {
    if (b < 16) {
        return "0" + b.toString(16);
    }
    return b.toString(16);
}

class Register {
    value;

    // Array of 8 div elements, each representing a single byte.
    // First element represents the most significant byte in register.
    bytes;

    /**
     * @param {string} name Register name.
    */
    constructor(name) {
        let r = document.getElementById(name);
        this.bytes = r.getElementsByClassName("byte");
    }

    /**
     * @param {Array<number>} value Array with 8 numbers. First element
     * corresponds to least significant byte in register.
    */
    changeValue(value) {
        this.value = value;
        
        let j = 8;
        for (let i = 0; i < 8; i += 1) {
            j -= 1;
            this.bytes[j].innerText = formatByte(value[i]);
        }
    }
}

main();
