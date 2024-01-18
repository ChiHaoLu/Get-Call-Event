var tracer = {
    data: [],
    fault: function (log) {
    },
    step: function (log) {
        var topicCount = (log.op.toString().match(/LOG(\d)/) || [])[1];
        if (topicCount) {
            var res = {
                address: log.contract.getAddress(),
                data: log.memory.slice(parseInt(log.stack.peek(0)), parseInt(log.stack.peek(0)) + parseInt(log.stack.peek(1))),
            };
            for (var i = 0; i < topicCount; i++)
                res[`topic` + i.toString()] = log.stack.peek(i + 2);
            this.data.push(res);
        }
    },
    result: function () {
        return this.data;
    }
}