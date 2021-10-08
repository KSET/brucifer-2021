export const state = () => ({});

export const getters = {};

export const mutations = {};

export const actions = {
  async nuxtServerInit({ dispatch }, context) {
    const getNuxtServerInits = (modules) => {
      const initFunctions = [];

      const descend = (modules, root = "") => {
        const moduleNames = Object.keys(modules || {});

        for (const moduleName of moduleNames) {
          const initFunctionName = `${ root }${ moduleName }/nuxtServerInit`;

          if (this._actions[initFunctionName]) {
            initFunctions.push(initFunctionName);
          }

          descend(modules[moduleName]._children, `${ root }${ moduleName }/`);
        }

        return modules;
      };

      descend(modules);

      return initFunctions.flat();
    };

    const nuxtServerInits = getNuxtServerInits(context.store._modules.root._children);

    await Promise.all(
      nuxtServerInits.map((initFnName) => dispatch(initFnName, context)),
    );
  },
};
